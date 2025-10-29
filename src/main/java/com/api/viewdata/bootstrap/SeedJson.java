package com.api.viewdata.bootstrap;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.api.viewdata.model.Status;
import com.api.viewdata.model.Transaction;
import com.api.viewdata.repo.StatusRepo;
import com.api.viewdata.repo.TransactionRepo;
import org.springframework.boot.CommandLineRunner;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

import java.io.InputStream;
import java.time.LocalDateTime;
import java.time.format.DateTimeFormatter;

@Configuration
public class SeedJson {
    private static final DateTimeFormatter FMT = DateTimeFormatter.ofPattern("yyyy-MM-dd HH:mm:ss");

    @Bean
    CommandLineRunner seed(ObjectMapper mapper, StatusRepo statusRepo, TransactionRepo txRepo) {
        return args -> {
            if (txRepo.count() > 0L) return;

            try (InputStream is = getClass().getResourceAsStream("/viewData.json")) {
                ViewDataFile file = mapper.readValue(is, ViewDataFile.class);

                for (ViewDataFile.StatusItem s : file.status) {
                    Status st = new Status();
                    st.setId(s.id);
                    st.setName(s.name);
                    statusRepo.save(st);
                }

                for (ViewDataFile.Row r : file.data) {
                    Transaction t = new Transaction();
                    t.setId(r.id);
                    t.setProductId(r.productID);
                    t.setProductName(r.productName);
                    t.setAmountText(r.amount);
                    t.setCustomerName(r.customerName);

                    Status st = statusRepo.findById(r.status)
                            .orElseThrow(() -> new IllegalStateException("Missing status id " + r.status));
                    t.setStatus(st);

                    t.setTransactionDate(LocalDateTime.parse(r.transactionDate, FMT));
                    t.setCreateBy(r.createBy);
                    t.setCreateOn(LocalDateTime.parse(r.createOn, FMT));

                    txRepo.save(t);
                }
            }
        };
    }

    @JsonIgnoreProperties(ignoreUnknown = true)
    static class ViewDataFile {
        public Row[] data;
        public StatusItem[] status;
        static class Row {
            public long id;
            public String productID;
            public String productName;
            public String amount;
            public String customerName;
            public int status;
            public String transactionDate;
            public String createBy;
            public String createOn;
        }
        static class StatusItem {
            public int id;
            public String name;
        }
    }
}