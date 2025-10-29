package com.api.viewdata.controller;

import org.springframework.web.bind.annotation.*;
import org.springframework.http.MediaType;
import com.api.viewdata.dto.ViewDataResponse;
import com.api.viewdata.service.ViewDataService;

@RestController
@RequestMapping("/api")
public class ViewDataController {
    private final ViewDataService service;
    public ViewDataController(ViewDataService service) { this.service = service; }

    @GetMapping(value="/view-data", produces=MediaType.APPLICATION_JSON_VALUE)
    public ViewDataResponse get() {
        return service.buildResponse();
    }
}