using Microsoft.AspNetCore.Authorization;
using Microsoft.AspNetCore.Mvc;
using Ports.InBound.Contracts;
using Ports.InBound.DTO;
using Ports.OutBound.DTO;
using Presentation.ResponseDTO.Analysis;

namespace Presentation.Controllers;

[ApiController]
[Route("api/analysis")]
public class AnalysisController : ControllerBase
{
    
    private readonly IAnalysisService _analysisService;
    
    private readonly IUserService _userService;

    public AnalysisController(IAnalysisService analysisService, IUserService userService)
    {
        _analysisService = analysisService;
        _userService = userService;
    }
    
    [Authorize(Roles = "Teacher")]
    [HttpPost("{labId}/run")]
    public ActionResult<PostAnalysisResponse> RunAnalysis([FromRoute] Guid labId)
    {
        var stringId = HttpContext.User.Claims
            .FirstOrDefault(c => c.Type == "UserId")!.Value;
        var userId = Guid.Parse(stringId);
        var runResult = _analysisService.StartAnalysisJob(labId, userId);
        if (runResult is AnalysisJobResult.Failure failure)
        {
            return BadRequest(failure.Reason);
        }
        
        var success = (AnalysisJobResult.Success)runResult;
        return Ok(new PostAnalysisResponse(
            success.AnalysisJob.JobId,
            success.AnalysisJob.LabId,
            success.AnalysisJob.Status));
    }

    [Authorize(Roles = "Teacher")]
    [HttpGet("{jobId}/jobs")]
    public ActionResult<JobStatusResponse> GetJobAnalysis([FromRoute] Guid jobId)
    {
        var stringId = HttpContext.User.Claims
            .FirstOrDefault(c => c.Type == "UserId")!.Value;
        var userId = Guid.Parse(stringId);
        
        
        var jobResult = _analysisService.GetAnalysisJobById(jobId, userId);
        if (jobResult is AnalysisJobResult.Failure jobFailure)
        {
            return BadRequest(jobFailure.Reason);
        }
        var success = (AnalysisJobResult.Success)jobResult;

        var teacherResult = _userService.GetUser(success.AnalysisJob.UserId);
        if (teacherResult is UserResult.Failure userFailure)
        {
            return BadRequest(userFailure.Reason);
        }
        var teacher = (UserResult.Success)teacherResult;
            
        return Ok(new JobStatusResponse(
            success.AnalysisJob.JobId,
            success.AnalysisJob.LabId,
            success.AnalysisJob.Status,
            new CreatedBy(teacher.User.UserId, teacher.User.FullName),
            success.AnalysisJob.CreatedAt,
            success.AnalysisJob.StartedAt,
            success.AnalysisJob.FinishedAt,
            success.AnalysisJob.ErrorMessage));

    }

    [Authorize(Roles = "Teacher")]
    [HttpGet("{labId}/labs")]
    public ActionResult<AnalysisResultResponse> GetAnalysis([FromRoute] Guid labId)
    {
        var stringId = HttpContext.User.Claims
            .FirstOrDefault(c => c.Type == "UserId")!.Value;
        var userId = Guid.Parse(stringId);
        
        var analysisResult = _analysisService.FullAnalysis(labId, userId);
        if (analysisResult is FullLabAnalysisResult.Failure failure)
        {
            return BadRequest(failure.Reason);
        }
        var success = (FullLabAnalysisResult.Success)analysisResult;

        return Ok(new AnalysisResultResponse(
            success.Job.LabId,
            new LastAnalysisJob(
                success.Job.JobId,
                success.Job.Status,
                success.Job.CreatedAt,
                success.Job.FinishedAt),
            success.Stats,
            success.SubItems
            ));
    }
    
}
