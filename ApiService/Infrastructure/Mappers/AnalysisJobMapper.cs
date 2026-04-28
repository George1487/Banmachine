using Domain.Jobs;
using Infrastructure.Entities;

namespace Infrastructure.Mappers;

public static class AnalysisJobMapper
{
    public static AnalysisJob ToDomain(AnalysisJobEntity entity)
    {
        var status = entity.Status switch
        {
            "pending" => JobStatus.Pending,
            "processing" => JobStatus.Processing,
            "done" => JobStatus.Done,
            _ => JobStatus.Failed
        };

        return new AnalysisJob(
            entity.Id,
            entity.LabId,
            status,
            entity.CreatedBy,
            entity.CreatedAt,
            entity.StartedAt,
            entity.FinishedAt,
            entity.ErrorMessage ?? string.Empty);
    }

    public static AnalysisJobEntity ToEntity(AnalysisJob domain)
    {
        var status = domain.Status switch
        {
            JobStatus.Pending => "pending",
            JobStatus.Processing => "processing",
            JobStatus.Done => "done",
            _ => "failed"
        };

        return new AnalysisJobEntity
        {
            Id = domain.JobId,
            LabId = domain.LabId,
            Status = status,
            CreatedBy = domain.UserId,
            CreatedAt = domain.CreatedAt,
            StartedAt = domain.StartedAt,
            FinishedAt = domain.FinishedAt,
            ErrorMessage = domain.ErrorMessage
        };
    }
}
