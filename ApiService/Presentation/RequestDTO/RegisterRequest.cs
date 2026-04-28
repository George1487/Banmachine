using Domain.Users;

namespace Presentation.RequestDTO;

public sealed record RegisterRequest(
    
    string Email,
    
    string FullName,
    
    string Password,
    
    string? GroupName,
    
    UserRole Role);