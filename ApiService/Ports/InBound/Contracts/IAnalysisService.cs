using Ports.InBound.DTO;
using Ports.OutBound.DTO;

namespace Ports.InBound.Contracts;

public interface IAnalysisService
{
    AnalysisJobResult StartAnalysisJob(Guid labId, Guid userId);
    
    AnalysisJobResult GetAnalysisJobById(Guid jobId, Guid userId);
    
    FullLabAnalysisResult FullAnalysis(Guid labId, Guid userId);
    
    OneSubmissionAnalysisResult GetOneSubmissionAnalysisResult(Guid submissionId,
        Guid userId);
}