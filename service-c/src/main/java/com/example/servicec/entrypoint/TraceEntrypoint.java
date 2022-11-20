package com.example.servicec.entrypoint;

import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RestController;

import io.micrometer.core.instrument.MeterRegistry;

@RestController
public class TraceEntrypoint {
	private final MeterRegistry registry;

	public TraceEntrypoint(final MeterRegistry registry) {
		this.registry = registry;
	}

	@GetMapping("/trace")
	public ResponseEntity<String> getTraceStep(){
		registry.counter("trace_step_total").increment();
		return ResponseEntity.ok("{\"massage\": \"trace\"}");
	}
}
