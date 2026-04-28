namespace Domain.Users;

public sealed record User(
    
    Guid UserId,
    
    string Email,
    
    string FullName,
    
    string Password,
    
    string? GroupName,
    
    UserRole Role);