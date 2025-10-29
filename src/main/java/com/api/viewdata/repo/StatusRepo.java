package com.api.viewdata.repo;

import org.springframework.data.jpa.repository.JpaRepository;
import com.api.viewdata.model.Status;

public interface StatusRepo extends JpaRepository<Status, Integer> {}