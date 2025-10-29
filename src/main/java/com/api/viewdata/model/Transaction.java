package com.api.viewdata.model;

import jakarta.persistence.*;
import java.time.LocalDateTime;

@Entity @Table(name="transactions")
public class Transaction {
    @Id
    private Long id;

    @Column(name="product_id", nullable=false)
    private String productId;

    @Column(name="product_name", nullable=false)
    private String productName;

    @Column(name="amount_text", nullable=false)
    private String amountText;

    @Column(name="customer_name", nullable=false)
    private String customerName;

    @ManyToOne(fetch = FetchType.LAZY)
    @JoinColumn(name="status_id", nullable=false)
    private Status status;

    @Column(name="transaction_date", nullable=false)
    private LocalDateTime transactionDate;

    @Column(name="create_by", nullable=false)
    private String createBy;

    @Column(name="create_on", nullable=false)
    private LocalDateTime createOn;

    public Long getId() { return id; }
    public void setId(Long id) { this.id = id; }
    public String getProductId() { return productId; }
    public void setProductId(String productId) { this.productId = productId; }
    public String getProductName() { return productName; }
    public void setProductName(String productName) { this.productName = productName; }
    public String getAmountText() { return amountText; }
    public void setAmountText(String amountText) { this.amountText = amountText; }
    public String getCustomerName() { return customerName; }
    public void setCustomerName(String customerName) { this.customerName = customerName; }
    public Status getStatus() { return status; }
    public void setStatus(Status status) { this.status = status; }
    public LocalDateTime getTransactionDate() { return transactionDate; }
    public void setTransactionDate(LocalDateTime transactionDate) { this.transactionDate = transactionDate; }
    public String getCreateBy() { return createBy; }
    public void setCreateBy(String createBy) { this.createBy = createBy; }
    public LocalDateTime getCreateOn() { return createOn; }
    public void setCreateOn(LocalDateTime createOn) { this.createOn = createOn; }
}