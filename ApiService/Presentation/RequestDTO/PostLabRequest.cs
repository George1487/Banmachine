namespace Presentation.RequestDTO;

public sealed record PostLabRequest(
    
    string Title,
    
    string Description,
    
    DateTimeOffset DeadlineAt
    
    ) { }