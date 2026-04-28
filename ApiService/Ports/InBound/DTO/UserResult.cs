using Domain.Users;

namespace Ports.InBound.DTO;

public abstract record UserResult
{
    
    private UserResult() { }
    
    public sealed record Success(User User) : UserResult();
    
    public sealed record Failure(string Reason) : UserResult();
}