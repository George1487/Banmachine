using Domain.Users;
using Infrastructure.Mappers;
using Microsoft.EntityFrameworkCore;
using Ports.InBound.DTO;

namespace Infrastructure.Repositories;

public class UserRepository
{
    private readonly AppDbContext _context;

    public UserRepository(AppDbContext context)
    {
        _context = context;
    }

    public UserResult AddUser(User user)
    {
        try
        {
            if (_context.Users.AsNoTracking().Any(x => x.Email == user.Email))
            {
                return new UserResult.Failure("user_already_exists");
            }

            var entity = UserMapper.ToEntity(user);
            _context.Users.Add(entity);
            _context.SaveChanges();

            return new UserResult.Success(UserMapper.ToDomain(entity));
        }
        catch (Exception exception)
        {
            return new UserResult.Failure(exception.Message);
        }
    }

    public UserResult GetUser(Guid userId)
    {
        try
        {
            var entity = _context.Users
                .AsNoTracking()
                .FirstOrDefault(x => x.Id == userId);

            return entity is null
                ? new UserResult.Failure("user_not_found")
                : new UserResult.Success(UserMapper.ToDomain(entity));
        }
        catch (Exception exception)
        {
            return new UserResult.Failure(exception.Message);
        }
    }

    public UserResult GetUserByEmail(string email)
    {
        try
        {
            var entity = _context.Users
                .AsNoTracking()
                .FirstOrDefault(x => x.Email == email);

            return entity is null
                ? new UserResult.Failure("user_not_found")
                : new UserResult.Success(UserMapper.ToDomain(entity));
        }
        catch (Exception exception)
        {
            return new UserResult.Failure(exception.Message);
        }
    }
}
