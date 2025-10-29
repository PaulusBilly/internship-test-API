package com.api.viewdata.dto;

import java.util.List;

public class ViewDataResponse {
    public List<Row> data;
    public List<StatusItem> status;

    public static class Row {
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

    public static class StatusItem {
        public int id;
        public String name;
    }
}