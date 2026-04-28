using Domain.Labs;
using Infrastructure.Entities;

namespace Infrastructure.Mappers;

public static class LabMapper
{
    public static Lab ToDomain(LabEntity entity)
    {
        return new Lab(
            entity.Id,
            entity.TeacherId,
            entity.Title,
            entity.Description ?? string.Empty,
            entity.Status == "active" ? LabStatus.Active : LabStatus.Closed,
            entity.DeadlineAt);
    }

    public static LabEntity ToEntity(Lab domain)
    {
        return new LabEntity
        {
            Id = domain.LabId,
            TeacherId = domain.TeacherId,
            Title = domain.Title,
            Description = domain.Description,
            Status = domain.LabStatus == LabStatus.Active ? "active" : "closed",
            DeadlineAt = domain.DeadlineAt
        };
    }
}
