package com.github.mankit007.recipe.Dashboard.Service.repository;

import org.springframework.data.domain.Pageable;
import org.springframework.data.mongodb.repository.ReactiveMongoRepository;
import org.springframework.stereotype.Repository;

import com.github.mankit007.recipe.Dashboard.Service.model.ClickEvent;

import reactor.core.publisher.Flux;
import reactor.core.publisher.Mono;

@Repository
public interface ClickEventRepository extends ReactiveMongoRepository<ClickEvent, String> {

    Flux<ClickEvent> findByTenantIdOrderByTimestampDesc(String tenantId, Pageable pageable);

    Mono<Long> countByTenantId(String tenantId);
}
