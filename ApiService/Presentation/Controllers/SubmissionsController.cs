using Domain.Users;
using Microsoft.AspNetCore.Authorization;
using Microsoft.AspNetCore.Mvc;
using Microsoft.AspNetCore.Http;
using Ports.InBound.Contracts;
using Ports.InBound.DTO;
using Ports.OutBound.Contracts;
using Ports.OutBound.DTO;
using Presentation.RequestDTO;
using Presentation.ResponseDTO.Submissions;

namespace Presentation.Controllers;

[ApiController]
[Route("api/submissions")]
public class SubmissionsController  : ControllerBase
{
    private readonly ISubmissionService _submissionService;
    
    private readonly ILabService _labService;
    
    private readonly IUserService _userService;
    
    private readonly IAnalysisService _analysisService;

    private readonly ISubmissionAnalysisSummaryService _summaryService;

    public SubmissionsController(
        ISubmissionService submissionService,
        ILabService labService,
        IUserService userService,
        IAnalysisService analysisService,
        ISubmissionAnalysisSummaryService summaryService)
    {
        _submissionService = submissionService;
        _labService = labService;
        _userService = userService;
        _analysisService = analysisService;
        _summaryService = summaryService;
    }

    [Authorize(Roles = "Student")]
    [HttpPost("{labId}")]
    public ActionResult<CreateSubmissionResponse> CreateSubmission(
        [FromRoute] Guid labId,
        [FromForm] CreateSubmissionRequest request)
    {
        var stringId = HttpContext.User.Claims
            .FirstOrDefault(c => c.Type == "UserId")!.Value;
        var userId = Guid.Parse(stringId);
        
        var result = _submissionService.CreateSubmission(request.File, labId, userId);
        if (result is SubmissionResult.Failure failure)
        {
            return BadRequest(failure.Reason);
        }
        
        var success = (SubmissionResult.Success)result;
        return Ok(new CreateSubmissionResponse(
            success.Submission.SubmissionId,
            success.Submission.LabId,
            success.Submission.StudentId,
            success.Submission.Status,
            success.Submission.SubmittedAt
        ));
    }

    [Authorize(Roles = "Student")]
    [HttpGet("me")]
    public ActionResult<List<GetSubmissionsResponse>> GetSubmissions()
    {
        var stringId = HttpContext.User.Claims
            .FirstOrDefault(c => c.Type == "UserId")!.Value;
        var userId = Guid.Parse(stringId);
        
        var subResult = _submissionService.GetSubmissionsByUserId(userId);
        if (subResult is SubmissionsResult.Failure subFailure)
        {
            return Ok(subFailure.Reason);
        }

        var subSuccess = (SubmissionsResult.Success)subResult;
        
        var submissions = new List<GetSubmissionsResponse>();
        foreach (var submission in subSuccess.Submissions)
        {
            var labInfo = _labService.GetLab(submission.LabId);
            if (labInfo is LabResult.Failure failure)
            {
                return BadRequest(failure.Reason);
            }
            var labSuccess = (LabResult.Success)labInfo;
            submissions.Add(new GetSubmissionsResponse(
                submission.SubmissionId,
                submission.LabId,
                labSuccess.Lab.Title,
                submission.Status,
                submission.SubmittedAt));
        }
        return Ok(submissions);
    }

    [Authorize()]
    [HttpGet("{submissionId}")]
    public ActionResult<GetConcreteSubmissionResponse> GetSubmission(
        [FromRoute] Guid submissionId)
    {
        var result = _submissionService.GetSubmissionById(submissionId);
        if (result is SubmissionResult.Failure subFailure)
        {
            return BadRequest(subFailure.Reason);
        }
        var subSuccess = (SubmissionResult.Success)result;
        return Ok(subSuccess.Submission);
    }

    [Authorize(Roles = "Teacher")]
    [HttpGet("labs/{labId}")]
    public ActionResult<List<GetSubmissionsTeacherResponse>> GetTeacherSubmissions(
        [FromRoute] Guid labId)
    {
        var stringId = HttpContext.User.Claims
            .FirstOrDefault(c => c.Type == "UserId")!.Value;
        var userId = Guid.Parse(stringId);
        
        
        var labInfo = _labService.GetLab(labId);
        if (labInfo is LabResult.Failure failure)
        {
            return BadRequest(failure.Reason);
        }
        var labSuccess = (LabResult.Success)labInfo;
        if (labSuccess.Lab.TeacherId != userId)
        {
            return BadRequest("You are not the owner of this lab");
        }
        
        
        var submissionsResult = _submissionService.GetSubmissionsByLabId(labId);
        if (submissionsResult is SubmissionsResult.Failure subFailure)
        {
            return BadRequest(subFailure.Reason);
        }
        var subSuccess = (SubmissionsResult.Success)submissionsResult;
        var teachersSubmissions = new List<GetSubmissionsTeacherResponse>();
        foreach (var submission in subSuccess.Submissions)
        {
            var userInfo = _userService.GetUser(submission.StudentId);
            if (userInfo is UserResult.Failure userFailure)
            {
                return BadRequest(userFailure.Reason);
            }
            var userSuccess = (UserResult.Success)userInfo;
            var studentInfo = new Students(userSuccess.User.UserId, userSuccess.User.FullName);

            decimal? topMatchScore = null;
            string? finalScoreRiskLevel = null;
            var summaryResult = _summaryService
                .GetSubmissionAnalysisSummaryBySubmissionId(submission.SubmissionId);
            if (summaryResult is SubmissionAnalysisSummaryResult.Success summarySuccess)
            {
                topMatchScore = summarySuccess.Summary.TopMatchScore;
                finalScoreRiskLevel = summarySuccess.Summary.FinalScoreRiskLevel;
            }

            teachersSubmissions.Add(new GetSubmissionsTeacherResponse(
                submission.SubmissionId,
                studentInfo,
                submission.Status,
                submission.SubmittedAt,
                topMatchScore,
                finalScoreRiskLevel));
        }
        return Ok(teachersSubmissions);
    }

    [Authorize(Roles = "Teacher")]
    [HttpGet("{submissionId}/matches")]
    public ActionResult<GetSubmissionsMatchesResponse> GetMatches(
        [FromRoute] Guid submissionId)
    {
        var stringId = HttpContext.User.Claims
            .FirstOrDefault(c => c.Type == "UserId")!.Value;
        var userId = Guid.Parse(stringId);

        var analysisResult = _analysisService.GetOneSubmissionAnalysisResult(submissionId, userId);
        if (analysisResult is OneSubmissionAnalysisResult.Failure failure)
        {
            return BadRequest(failure.Reason);
        }
        var success = (OneSubmissionAnalysisResult.Success)analysisResult;

        return Ok(new GetSubmissionsMatchesResponse(
            success.SubmissionId,
            success.JobId,
            success.SubMatches));
    }
}
