using Domain.Submissions;
using Infrastructure.Entities;

namespace Infrastructure.Mappers;

public static class SubmissionMapper
{
    public static Submission ToDomain(SubmissionEntity entity)
    {
        var status = entity.Status switch
        {
            "uploaded" => SubmissionStatus.Uploaded,
            "parsed" => SubmissionStatus.Parsed,
            "parsing" => SubmissionStatus.Parsing,
            _ => SubmissionStatus.Failed
        };

        return new Submission(
            entity.Id,
            entity.LabId,
            entity.StudentId,
            status,
            entity.MimeType,
            entity.SourceFileName,
            entity.StorageKey,
            entity.SubmittedAt
        );
    }

    public static SubmissionEntity ToEntity(Submission domain)
    {
        var status = domain.Status switch
        {
            SubmissionStatus.Uploaded => "uploaded",
            SubmissionStatus.Parsed => "parsed",
            SubmissionStatus.Parsing => "parsing",
            _ => "failed"
        };

        return new SubmissionEntity
        {
            Id = domain.SubmissionId,
            LabId = domain.LabId,
            StudentId = domain.StudentId,
            SubmittedAt = domain.SubmittedAt,
            MimeType = domain.MimeType,
            SourceFileName = domain.SourceFileName,
            StorageKey = domain.StorageKey,
            Status = status
        };
    }
}
