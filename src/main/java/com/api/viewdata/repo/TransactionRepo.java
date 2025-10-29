package com.api.viewdata.repo;

import org.springframework.data.jpa.repository.JpaRepository;
import com.api.viewdata.model.Transaction;

public interface TransactionRepo extends JpaRepository<Transaction, Long> {}