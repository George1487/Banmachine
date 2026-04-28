using Domain.Labs;
using Infrastructure.Mappers;
using Microsoft.EntityFrameworkCore;
using Ports.InBound.DTO;

namespace Infrastructure.Repositories;

public class LabRepository
{
    private readonly AppDbContext _context;

    public LabRepository(AppDbContext context)
    {
        _context = context;
    }

    public LabResult GetLab(Guid labId)
    {
        try
        {
            var entity = _context.Labs
                .AsNoTracking()
                .FirstOrDefault(x => x.Id == labId);

            return entity is null
                ? new LabResult.Failure("lab_not_found")
                : new LabResult.Success(LabMapper.ToDomain(entity));
        }
        catch (Exception exception)
        {
            return new LabResult.Failure(exception.Message);
        }
    }

    public LabResult GetLabByTitle(string title)
    {
        try
        {
            var entity = _context.Labs
                .AsNoTracking()
                .FirstOrDefault(x => x.Title == title);

            return entity is null
                ? new LabResult.Failure("lab_not_found")
                : new LabResult.Success(LabMapper.ToDomain(entity));
        }
        catch (Exception exception)
        {
            return new LabResult.Failure(exception.Message);
        }
    }

    public LabsResult GetLabs()
    {
        try
        {
            var labs = _context.Labs
                .AsNoTracking()
                .OrderBy(x => x.DeadlineAt)
                .Select(LabMapper.ToDomain)
                .ToList();

            return new LabsResult.Success(labs);
        }
        catch (Exception exception)
        {
            return new LabsResult.Failure(exception.Message);
        }
    }

    public LabResult AddLab(Lab lab)
    {
        try
        {
            if (_context.Labs.AsNoTracking().Any(x => x.Title == lab.Title))
            {
                return new LabResult.Failure("lab_already_exists");
            }

            var entity = LabMapper.ToEntity(lab);
            _context.Labs.Add(entity);
            _context.SaveChanges();

            return new LabResult.Success(LabMapper.ToDomain(entity));
        }
        catch (Exception exception)
        {
            return new LabResult.Failure(exception.Message);
        }
    }

    public LabResult PatchLab(
        Guid labId,
        LabStatus newLabStatus,
        string newTitle,
        string newDescription,
        DateTimeOffset newDeadlineAt)
    {
        try
        {
            var entity = _context.Labs.FirstOrDefault(x => x.Id == labId);
            if (entity is null)
            {
                return new LabResult.Failure("lab_not_found");
            }

            var titleInUse = _context.Labs
                .AsNoTracking()
                .Any(x => x.Title == newTitle && x.Id != labId);
            if (titleInUse)
            {
                return new LabResult.Failure("lab_already_exists");
            }

            entity.Title = newTitle;
            entity.Description = newDescription;
            entity.Status = newLabStatus == LabStatus.Active ? "active" : "closed";
            entity.DeadlineAt = newDeadlineAt;

            _context.SaveChanges();

            return new LabResult.Success(LabMapper.ToDomain(entity));
        }
        catch (Exception exception)
        {
            return new LabResult.Failure(exception.Message);
        }
    }
}
