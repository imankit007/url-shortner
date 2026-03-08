package com.github.mankit007.recipe.Dashboard.Service.repository;

import org.springframework.data.mongodb.repository.ReactiveMongoRepository;
import org.springframework.stereotype.Repository;

import com.github.mankit007.recipe.Dashboard.Service.model.ClickAggregate;

import reactor.core.publisher.Flux;

@Repository
public interface ClickAggregateRepository extends ReactiveMongoRepository<ClickAggregate, String> {

    Flux<ClickAggregate> findByTenantIdAndGranularityOrderByWindowStartDesc(
            String tenantId, String granularity);

    Flux<ClickAggregate> findByTenantIdAndShortCodeAndGranularityOrderByWindowStartAsc(
            String tenantId, String shortCode, String granularity);
}
