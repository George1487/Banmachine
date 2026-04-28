using Domain.Jobs;
using Ports.OutBound.DTO;

namespace Ports.OutBound.Contracts;

public interface IAnalysisRepository
{
    
    AnalysisJobResult GetAnalysisJobById(Guid jobId);
    
    AnalysisJobResult AddAnalysisJob(AnalysisJob analysisJob);
    
    AnalysisJobResult GetAnalysisJobByLabId(Guid labId);
    
    AnalysisJobResult GetAnalysisJobBySubmissionId(Guid submissionId);
    
    AnalysisJobResult GetAnalysisJobByLastDate();
    
    AnalysisStatsResults GetAnalysisStatsByLabId(Guid labId);
}