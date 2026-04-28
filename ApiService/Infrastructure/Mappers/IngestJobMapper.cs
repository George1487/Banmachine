using Domain.Jobs;
using Infrastructure.Entities;

namespace Infrastructure.Mappers;

public static class IngestJobMapper
{
    public static IngestJob ToDomain(IngestJobEntity entity)
    {
        var status = entity.Status switch
        {
            "pending" => JobStatus.Pending,
            "processing" => JobStatus.Processing,
            "done" => JobStatus.Done,
            _ => JobStatus.Failed
        };

        return new IngestJob(
            entity.Id,
            entity.SubmissionId,
            status,
            entity.CreatedAt,
            null,
            entity.FinishedAt,
            entity.ErrorMessage ?? string.Empty);
    }

    public static IngestJobEntity ToEntity(IngestJob domain)
    {
        var status = domain.Status switch
        {
            JobStatus.Pending => "pending",
            JobStatus.Processing => "processing",
            JobStatus.Done => "done",
            _ => "failed"
        };

        return new IngestJobEntity
        {
            Id = domain.IngestJobId,
            SubmissionId = domain.SubmissionId,
            Status = status,
            CreatedAt = domain.CreatedAt,
            FinishedAt = domain.FinishedAt,
            ErrorMessage = domain.ErrorMessage
        };
    }
}
