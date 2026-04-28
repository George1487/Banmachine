namespace Ports.InBound.DTO;

public sealed record Student(
    Guid StudentId,
    string FullName);