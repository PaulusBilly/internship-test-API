package com.api.viewdata.service;

import org.springframework.stereotype.Service;
import com.api.viewdata.dto.ViewDataResponse;
import com.api.viewdata.model.Transaction;
import com.api.viewdata.model.Status;
import com.api.viewdata.repo.TransactionRepo;
import com.api.viewdata.repo.StatusRepo;

import java.time.format.DateTimeFormatter;
import java.util.ArrayList;

@Service
public class ViewDataService {
    private final TransactionRepo txRepo;
    private final StatusRepo statusRepo;
    private static final DateTimeFormatter FMT = DateTimeFormatter.ofPattern("yyyy-MM-dd HH:mm:ss");

    public ViewDataService(TransactionRepo txRepo, StatusRepo statusRepo) {
        this.txRepo = txRepo;
        this.statusRepo = statusRepo;
    }

    public ViewDataResponse buildResponse() {
        var res = new ViewDataResponse();

        res.data = new ArrayList<>();
        for (Transaction t : txRepo.findAll()) {
            var row = new ViewDataResponse.Row();
            row.id = t.getId();
            row.productID = t.getProductId();          // exact field name
            row.productName = t.getProductName();
            row.amount = t.getAmountText();            // string
            row.customerName = t.getCustomerName();
            row.status = t.getStatus().getId();
            row.transactionDate = t.getTransactionDate().format(FMT);
            row.createBy = t.getCreateBy();
            row.createOn = t.getCreateOn().format(FMT);
            res.data.add(row);
        }

        res.status = new ArrayList<>();
        for (Status s : statusRepo.findAll()) {
            var si = new ViewDataResponse.StatusItem();
            si.id = s.getId();
            si.name = s.getName();
            res.status.add(si);
        }
        return res;
    }
}