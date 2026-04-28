namespace Domain.Submissions;

public sealed record Submission(

    Guid SubmissionId,

    Guid LabId,

    Guid StudentId,

    SubmissionStatus Status,

    string MimeType,

    string SourceFileName,

    string StorageKey,
    
    DateTimeOffset SubmittedAt
);   