using Domain.Jobs;
using Domain.Submissions;
using Ports.InBound.Contracts;
using Ports.InBound.DTO;
using Ports.OutBound.Contracts;
using Ports.OutBound.DTO;

namespace Application.Services;

public class AnalysisJobServiceImpl : IAnalysisService
{
    
    private readonly IAnalysisRepository _analysisRepository;
    
    private readonly ISubmissionService _submissionService;
    
    private readonly ILabService _labService;
    
    private readonly ISubmissionAnalysisSummaryService _summaryService;
    
    private readonly IUserService _userService;
    
    private IPairwiseSimilarityService _pairwiseSimilarityService;

    public AnalysisJobServiceImpl(
        IAnalysisRepository analysisRepository,
        ISubmissionService submissionService,
        ILabService labService,
        ISubmissionAnalysisSummaryService summaryService,
        IUserService userService,
        IPairwiseSimilarityService pairwiseSimilarityService)
    {
        _analysisRepository = analysisRepository;
        _submissionService = submissionService;
        _labService = labService;
        _summaryService = summaryService;
        _userService = userService;
        _pairwiseSimilarityService = pairwiseSimilarityService;
    }
    
    public AnalysisJobResult StartAnalysisJob(Guid labId, Guid userId)
    {
      var jobResult =  _analysisRepository.GetAnalysisJobByLabId(labId);
      if (jobResult is AnalysisJobResult.Success success && 
          (success.AnalysisJob.Status == JobStatus.Pending 
           || success.AnalysisJob.Status == JobStatus.Processing ))
      {
          return new AnalysisJobResult.Failure("conflict");
      }
      
      var parsedSubmissionsResult = _submissionService.GetParsedSubmissionsByLabId(labId);
      if (parsedSubmissionsResult is ParsedSubmissionsResult.Failure failure)
      {
          return new AnalysisJobResult.Failure(failure.Reason);
      }

      var parsedSubmissionsSuccess = (ParsedSubmissionsResult.Success)parsedSubmissionsResult;
      var parsedSubmissionsCount = parsedSubmissionsSuccess
          .ParsedSubmission.Count;
      if (parsedSubmissionsCount < 2)
      {
          var oneAnalysisJob = new AnalysisJob(
              Guid.NewGuid(),
              labId,
              JobStatus.Done,
              userId,
              DateTimeOffset.Now,
              null,
              null,
              ""
          );
          return _analysisRepository.AddAnalysisJob(oneAnalysisJob);
      }

      var analysisJob = new AnalysisJob(
          Guid.NewGuid(),
          labId,
          JobStatus.Pending,
          userId,
          DateTimeOffset.Now,
          null,
          null,
          "");
      return _analysisRepository.AddAnalysisJob(analysisJob);
    }

    public AnalysisJobResult GetAnalysisJobById(Guid jobId, Guid userId)
    {
        var jobResult = _analysisRepository.GetAnalysisJobById(jobId);
        if (jobResult is AnalysisJobResult.Failure failure)
        {
            return new AnalysisJobResult.Failure(failure.Reason);
        }
        var jobSuccess = (AnalysisJobResult.Success)jobResult;
        if (jobSuccess.AnalysisJob.UserId != userId)
        {
            return new AnalysisJobResult.Failure("you do not have sufficient rights");
        }
        return jobSuccess;
    }

    public FullLabAnalysisResult FullAnalysis(Guid labId, Guid userId)
    {
        var labResult = _labService.GetLab(labId);
        if (labResult is LabResult.Failure labFailure)
        {
            return new FullLabAnalysisResult.Failure(labFailure.Reason);
        }

        var successLab = (LabResult.Success)labResult;
        if (userId != successLab.Lab.TeacherId)
        {
            return new FullLabAnalysisResult.Failure("you do not have sufficient rights");
        }


        var lastJobResult = _analysisRepository.GetAnalysisJobByLabId(labId);
        if (lastJobResult is AnalysisJobResult.Failure jobFailure){
            return new FullLabAnalysisResult.Failure(jobFailure.Reason);
        }
        var successLastJob = (AnalysisJobResult.Success)lastJobResult;
        
        var statsResult = _analysisRepository.GetAnalysisStatsByLabId(labId);
        if (statsResult is AnalysisStatsResults.Failure statsFailure)
        {
            return new FullLabAnalysisResult.Failure(statsFailure.Reason);
        }
        var stats = (AnalysisStatsResults.Success)statsResult;
        
        var submissionsResult = _submissionService.GetSubmissionsByLabId(labId);
        if (submissionsResult is SubmissionsResult.Failure subFailure)
        {
            return new FullLabAnalysisResult.Failure(subFailure.Reason);
        }
        
        var successSubmissions = (SubmissionsResult.Success)submissionsResult;
        var subItems = new List<SubItem>();
        foreach (var submission in successSubmissions.Submissions)
        {
            var summaryResult = _summaryService
                .GetSubmissionAnalysisSummaryBySubmissionId(submission.SubmissionId);
            if (summaryResult is SubmissionAnalysisSummaryResult.Failure sumFailure)
            {
                if (sumFailure.Reason == "submission_analysis_summary_not_found")
                {
                    continue;
                }

                return new FullLabAnalysisResult.Failure(sumFailure.Reason);
            }
            var summary = (SubmissionAnalysisSummaryResult.Success)summaryResult;
            
            var userResult = _userService.GetUser(submission.StudentId);
            if (userResult is UserResult.Failure userFailure)
            {
                return new FullLabAnalysisResult.Failure(userFailure.Reason);
            }
            var user = (UserResult.Success)userResult;

            var subItem = new SubItem(
                submission.SubmissionId,
                new Student(user.User.UserId, user.User.FullName),
                summary.Summary.TopMatchSubmissionId,
                summary.Summary.TopMatchScore,
                summary.Summary.FinalScoreRiskLevel
                );
            
            subItems.Add(subItem);
        }

        return new FullLabAnalysisResult.Success(successLastJob.AnalysisJob, 
            stats.Stats, subItems);
    }

    public OneSubmissionAnalysisResult GetOneSubmissionAnalysisResult(Guid submissionId,
        Guid userId)
    {
        var submissionResult = _submissionService.GetSubmissionById(submissionId);
        if (submissionResult is SubmissionResult.Failure subFailure)
        {
            return new OneSubmissionAnalysisResult.Failure(subFailure.Reason);
        }
        var submission = (SubmissionResult.Success)submissionResult;

        var jobResult = _analysisRepository.GetAnalysisJobBySubmissionId(submissionId);
        if (jobResult is AnalysisJobResult.Failure jobFailure)
        {
            return new OneSubmissionAnalysisResult.Failure(jobFailure.Reason);
        }
        var job = (AnalysisJobResult.Success)jobResult;

        var pairwiseSimilarityResult = _pairwiseSimilarityService
            .GetPairwiseSimilarityBySubmissionId(submission.Submission.SubmissionId);
        if (pairwiseSimilarityResult is PairwiseSimilaritiesResult.Failure pairFailure)
        {
            return new OneSubmissionAnalysisResult.Failure(pairFailure.Reason);
        }
        var pairWiseSimilarity = (PairwiseSimilaritiesResult.Success)pairwiseSimilarityResult;
        
        var subMatches = new List<SubMatch>();
        foreach (var pairWise in pairWiseSimilarity.PairwiseSimilarity)
        {
            var matchedSubmissionResult = _submissionService.GetSubmissionById(pairWise.RightSubmissionId);
            if (matchedSubmissionResult is SubmissionResult.Failure matchedSubmissionFailure)
            {
                return new OneSubmissionAnalysisResult.Failure(matchedSubmissionFailure.Reason);
            }
            var matchedSubmission = (SubmissionResult.Success)matchedSubmissionResult;

            var userResult = _userService.GetUser(matchedSubmission.Submission.StudentId);
            if (userResult is UserResult.Failure userFailure)
            {
                return new OneSubmissionAnalysisResult.Failure(userFailure.Reason);
            }
            var user = (UserResult.Success)userResult;

            var summaryResult = _summaryService
                .GetSubmissionAnalysisSummaryBySubmissionId(pairWise.RightSubmissionId);
            if (summaryResult is SubmissionAnalysisSummaryResult.Failure summaryFailure)
            {
                return new OneSubmissionAnalysisResult.Failure(summaryFailure.Reason);
            }
            var summary = (SubmissionAnalysisSummaryResult.Success)summaryResult;

            var subMatch = new SubMatch(
                pairWise.RightSubmissionId,
                new Student(user.User.UserId, user.User.FullName),
                pairWise.TextScore,
                pairWise.CalculationScore,
                pairWise.ImagesScore,
                pairWise.FinalScore,
                summary.Summary.FinalScoreRiskLevel);
            subMatches.Add(subMatch);
        }
        
        return new OneSubmissionAnalysisResult.Success(submission.Submission.SubmissionId,
            job.AnalysisJob.JobId, subMatches);
    }
}
