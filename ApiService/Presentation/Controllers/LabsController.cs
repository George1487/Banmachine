using System.Security.Claims;
using Domain.Labs;
using Domain.Users;
using Microsoft.AspNetCore.Authorization;
using Microsoft.AspNetCore.Mvc;
using Ports.InBound.Contracts;
using Ports.InBound.DTO;
using Ports.OutBound.Contracts;
using Presentation.RequestDTO;
using Presentation.ResponseDTO;
using Presentation.ResponseDTO.Labs;
using ILabService = Ports.InBound.Contracts.ILabService;

namespace Presentation.Controllers;

[ApiController]
[Route("api/labs")]
public class LabsController : ControllerBase
{

    private readonly ILabService _labService;
    
    private readonly ISubmissionService _submissionService;

    public LabsController(ILabService labService, ISubmissionService submissionService)
    {
        _labService = labService;
        _submissionService = submissionService;
    }    
    
    [Authorize]
    [HttpGet()]
    public ActionResult<List<LabDefaultInfoResponse>> GetLabs()
    {
        var stringId = HttpContext.User.Claims
            .FirstOrDefault(c => c.Type == "UserId")!.Value;
        var userId = Guid.Parse(stringId);
        
        var roleString = HttpContext.User.Claims
            .FirstOrDefault(c => c.Type == ClaimTypes.Role)!.Value;
        var role = UserRole.Student;
        if (roleString == "Teacher")
        {
            role = UserRole.Teacher;
        }
        
        var result =  _labService.GetLabs(userId, role);
        if (result is LabsResult.Failure failure)
        {
            return BadRequest(failure.Reason);
        }
        
        var success = (LabsResult.Success)result;
        List<LabDefaultInfoResponse> labs = success
            .Labs
            .Select(lab => new LabDefaultInfoResponse(
                    lab.LabId,
                    lab.Title,
                    lab.LabStatus,
                    lab.DeadlineAt)
            )
            .ToList();
        return Ok(labs);
    }
    
    [HttpGet("{labId}")]
    [Authorize()]
    public ActionResult<LabDetailedInfoResponse> GetLab([FromRoute] Guid labId)
    {
        var result = _labService.GetLab(labId);
        if (result is LabResult.Failure failure)
        {
            return BadRequest(failure.Reason);
        }
        
        var success = (LabResult.Success)result;
        return Ok(new LabDetailedInfoResponse(
            success.Lab.LabId,
            success.Lab.Title,
            success.Lab.Description,
            success.Lab.LabStatus,
            success.Lab.DeadlineAt));
    }

    [Authorize(Roles = "Teacher")]
    [HttpGet("me")]
    public ActionResult<List<LabTeachersInfoResponse>> GetTeachers()
    {
        var stringId = HttpContext.User.Claims
            .FirstOrDefault(c => c.Type == "UserId")!.Value;
        var userId = Guid.Parse(stringId);
        
        var result =  _labService.GetLabs(userId, UserRole.Teacher);
        if (result is LabsResult.Failure failure)
        {
            return BadRequest(failure.Reason);
        }
        var success = (LabsResult.Success)result;
        
        var answer = new List<LabTeachersInfoResponse>();
        foreach (var lab in success.Labs)
        {
            var infoResult = _submissionService.GetLabSubmissionsInfo(lab.LabId);
            if (infoResult is LabSubmissionsInfoResult.Failure subFailure)
            {
                return BadRequest(subFailure.Reason);
            } 
            var subSuccess = (LabSubmissionsInfoResult.Success)infoResult;
            answer.Add(new LabTeachersInfoResponse(
                lab.LabId,
                lab.Title,
                lab.LabStatus,
                subSuccess.SubmissionCount,
                subSuccess.ParsedSubmissionCount
            ));
        }

        return Ok(answer);
    }
    
    
    [Authorize(Roles = "Teacher")]
    [HttpPost()]
    public ActionResult<LabDetailedInfoResponse> CreateLab(
        [FromBody] PostLabRequest request)
    {
        
        var stringId = HttpContext.User.Claims
            .FirstOrDefault(c => c.Type == "UserId")!.Value;
        
        var userId = Guid.Parse(stringId);

        var result = _labService.AddLab(new Lab(
            Guid.NewGuid(),
            userId,
            request.Title,
            request.Description,
            LabStatus.Active,
            request.DeadlineAt));

        if (result is LabResult.Failure failure)
        {
            return BadRequest(failure.Reason);
        }
        var success = (LabResult.Success)result;
        return Ok(new LabDetailedInfoResponse(
            success.Lab.LabId,
            success.Lab.Title,
            success.Lab.Description,
            success.Lab.LabStatus,
            success.Lab.DeadlineAt));
    }

    [Authorize(Roles = "Teacher")]
    [HttpPatch("{labId}")]
    public ActionResult<LabDetailedInfoResponse> PatchLab([FromRoute] Guid labId,
        [FromBody] PatchLabRequest request)
    {
        var stringId = HttpContext.User.Claims
            .FirstOrDefault(c => c.Type == "UserId")!.Value;
        var userId = Guid.Parse(stringId);
        
        var result = _labService.PatchLab( 
            request.LabId,
            userId,
            request.Status,
            request.Title,
            request.Description,
            request.DeadlineAt);
        if (result is LabResult.Failure failure)
        {
            return BadRequest(failure.Reason);
        }
        
        var success = (LabResult.Success)result;
        return Ok(new LabDetailedInfoResponse(success.Lab.LabId,
            success.Lab.Title,
            success.Lab.Description,
            success.Lab.LabStatus,
            success.Lab.DeadlineAt));
    }
}