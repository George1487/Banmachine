using Domain.Submissions;
using Ports.OutBound.DTO;

namespace Ports.OutBound.Contracts;

public interface ISubmissionRepository
{
    SubmissionResult GetSubmission(Guid submissionId);
    
    SubmissionsResult GetSubmissionsByLabId(Guid labId);
    
    SubmissionsResult GetSubmissions();
    
    SubmissionsResult GetSubmissionsByUserId(Guid userId);
    
    ParsedSubmissionsResult GetParsedSubmissions();
    
    SubmissionResult AddSubmission(Submission submission);
    
    ParsedSubmissionsResult GetParsedSubmissionByLabId(Guid labId);
    
}