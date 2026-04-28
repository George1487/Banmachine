using Microsoft.AspNetCore.Http;
using Ports.InBound.DTO;
using Ports.OutBound.DTO;

namespace Ports.InBound.Contracts;

public interface ISubmissionService
{
    LabSubmissionsInfoResult GetLabSubmissionsInfo(Guid labId);
    
    SubmissionResult CreateSubmission(
        IFormFile file, 
        Guid labId,
        Guid studentId);
    
    SubmissionsResult GetSubmissionsByUserId(Guid userId);
    
    SubmissionsResult GetSubmissionsByLabId(Guid labId);
    
    SubmissionResult GetSubmissionById(Guid submissionId);
    
    ParsedSubmissionsResult GetParsedSubmissionsByLabId(Guid labId);
}