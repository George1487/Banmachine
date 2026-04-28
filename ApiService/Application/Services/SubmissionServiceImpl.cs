using Domain.Jobs;
using Domain.Submissions;
using Microsoft.AspNetCore.Http;
using Microsoft.Extensions.Options;
using Minio;
using Minio.DataModel.Args;
using Ports.InBound.Contracts;
using Ports.InBound.DTO;
using Ports.OutBound.Contracts;
using Ports.OutBound.DTO;
using ILabService = Ports.InBound.Contracts.ILabService;

namespace Application.Services;

public class SubmissionServiceImpl : ISubmissionService
{
    
    private readonly ISubmissionRepository _submissionRepository;
    
    private readonly ILabService _labService;
    
    private readonly IMinioClient  _minioClient;
    
    private readonly IIngestJobRepository _ingestJobRepository;

    private readonly string _bucketName;

    public SubmissionServiceImpl(
        ISubmissionRepository submissionRepository,
        ILabService labService, 
        IMinioClient minioClient,
        IIngestJobRepository ingestJobRepository,
        IOptions<StorageSettings> storageOptions)
    {
        _submissionRepository = submissionRepository;
        _labService = labService;
        _minioClient = minioClient;
        _ingestJobRepository = ingestJobRepository;
        _bucketName = storageOptions.Value.Bucket;
    }
    
    public LabSubmissionsInfoResult GetLabSubmissionsInfo(Guid labId)
    {
        var submissionsResult = _submissionRepository
            .GetSubmissionsByLabId(labId);
        if (submissionsResult is SubmissionsResult.Failure failure)
        {
            return new LabSubmissionsInfoResult.Failure(failure.Reason);
        }
        
        var submissionsSuccess = (SubmissionsResult.Success)submissionsResult;
        var submissionsCount = submissionsSuccess
            .Submissions
            .Count;
        
        var submissionsIds = submissionsSuccess
            .Submissions
            .Select(o => o.SubmissionId)
            .ToList();

        var parsedSubmissionsResult = _submissionRepository.GetParsedSubmissions();
        if (parsedSubmissionsResult is ParsedSubmissionsResult.Failure parseFailure)
        {
           return new LabSubmissionsInfoResult.Failure(parseFailure.Reason); 
        }
        
        var parsedSubmissionsSuccess = (ParsedSubmissionsResult.Success)parsedSubmissionsResult;
        var parsedSubmissionsCount = parsedSubmissionsSuccess
            .ParsedSubmission
            .Where(o => submissionsIds.Contains(o.SubmissionId))
            .ToList()
            .Count;
        
        return new LabSubmissionsInfoResult.Success(submissionsCount, parsedSubmissionsCount);
    }

    public SubmissionResult CreateSubmission(
        IFormFile file,
        Guid labId,
        Guid studentId)
    {
        var labResult = _labService.GetLab(labId);
        if (labResult is LabResult.Failure failure)
        {
            return new SubmissionResult.Failure(failure.Reason);
        }
        var lab = (LabResult.Success)labResult;
        if (lab.Lab.DeadlineAt.ToUniversalTime() < DateTime.UtcNow)
        {
            return new SubmissionResult.Failure("Too late");
        }
        
        var mimeType = file.ContentType;
        var name = file.FileName;
        var extension = Path.GetExtension(name)?.ToLowerInvariant();
        if (extension != ".docx")
        {
            return new SubmissionResult.Failure("unsupported_file_type");
        }
        if (!string.Equals(
                mimeType,
                "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
                StringComparison.OrdinalIgnoreCase))
        {
            return new SubmissionResult.Failure("unsupported_file_type");
        }
        
        
        var key = $"{labId}/{Guid.NewGuid()}";
        var submissionObject = new PutObjectArgs()
            .WithObject(key)
            .WithStreamData(file.OpenReadStream())
            .WithBucket(_bucketName)
            .WithContentType(mimeType)
            .WithObjectSize(file.Length);
        try
        {
            _minioClient.PutObjectAsync(submissionObject).GetAwaiter().GetResult();
        }
        catch (Exception e)
        {
            return new SubmissionResult.Failure($"minioError:{e.Message}");
        }
        
        var submission = new Submission(
            Guid.NewGuid(),
            labId,
            studentId,
            SubmissionStatus.Uploaded,
            mimeType,
            name,
            key,
            DateTimeOffset.UtcNow
        );

        var createSubmissionResult = _submissionRepository.AddSubmission(submission);
        if (createSubmissionResult is SubmissionResult.Failure submissionFailure)
        {
            return new SubmissionResult.Failure(submissionFailure.Reason);
        }

        var createJobResult = _ingestJobRepository.AddIngestJob(
            new IngestJob(Guid.NewGuid(), 
                submission.SubmissionId,
                JobStatus.Pending,
                DateTimeOffset.UtcNow,
                null,
                null,
                ""));
        if (createJobResult is IngestJobResult.Failure jobFailure)
        {
            return new SubmissionResult.Failure(jobFailure.Reason);
        }

        return createSubmissionResult;
    }

    public SubmissionsResult GetSubmissionsByUserId(Guid userId)
    {
        return _submissionRepository.GetSubmissionsByUserId(userId);
    }

    public SubmissionsResult GetSubmissionsByLabId(Guid labId)
    {
        return _submissionRepository.GetSubmissionsByLabId(labId);
    }

    public SubmissionResult GetSubmissionById(Guid submissionId)
    {
        return _submissionRepository.GetSubmission(submissionId);
    }

    public ParsedSubmissionsResult GetParsedSubmissionsByLabId(Guid labId)
    {
        return _submissionRepository.GetParsedSubmissionByLabId(labId);
    }
}
