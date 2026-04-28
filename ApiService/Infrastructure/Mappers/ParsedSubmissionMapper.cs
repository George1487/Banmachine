using Domain.Submissions;
using Infrastructure.Entities;

namespace Infrastructure.Mappers;

public static class ParsedSubmissionMapper
{
    public static ParsedSubmission ToDomain(ParsedSubmissionEntity entity)
    {
        return new ParsedSubmission(
            entity.Id,
            entity.SubmissionId,
            entity.ParsedAt);
    }

    public static ParsedSubmissionEntity ToEntity(ParsedSubmission parsedSubmission)
    {
        return new ParsedSubmissionEntity
        {
            Id = parsedSubmission.ParsedSubmissionId,
            SubmissionId = parsedSubmission.SubmissionId,
            ParsedAt = parsedSubmission.ParsedAt,
            RawText = string.Empty,
            StructuredData = "{}"
        };
    }
}
