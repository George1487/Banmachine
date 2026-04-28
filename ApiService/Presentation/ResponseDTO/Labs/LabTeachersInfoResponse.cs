using Domain.Labs;

namespace Presentation.ResponseDTO.Labs;

public record LabTeachersInfoResponse(
    Guid LabId,
    string Title,
    LabStatus LabStatus,
    int SubmissionCount,
    int ParsedSubmissionCount);