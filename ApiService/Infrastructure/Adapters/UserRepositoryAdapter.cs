using Domain.Users;
using Infrastructure.Repositories;
using Ports.InBound.DTO;
using Ports.OutBound.Contracts;

namespace Infrastructure.Adapters;

public class UserRepositoryAdapter : IUserRepository
{
    private readonly UserRepository _repository;

    public UserRepositoryAdapter(UserRepository repository)
    {
        _repository = repository;
    }

    public UserResult AddUser(User user)
    {
        return _repository.AddUser(user);
    }

    public UserResult GetUser(Guid userId)
    {
        return _repository.GetUser(userId);
    }

    public UserResult GetUserByEmail(string email)
    {
        return _repository.GetUserByEmail(email);
    }
}
