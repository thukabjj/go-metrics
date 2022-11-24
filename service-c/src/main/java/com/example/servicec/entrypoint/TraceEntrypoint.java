package com.example.servicec.entrypoint;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.cloud.sleuth.Span;
import org.springframework.cloud.sleuth.Tracer;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RestController;

import com.example.servicec.entrypoint.entity.TraceEntrypointResponse;
import com.example.servicec.service.TraceService;
import io.micrometer.core.instrument.MeterRegistry;

@RestController
public class TraceEntrypoint {
	private final MeterRegistry registry;

	private final TraceService traceService;

	@Autowired
	private Tracer tracer;
	public TraceEntrypoint(final MeterRegistry registry, final TraceService traceService) {
		this.registry = registry;
		this.traceService = traceService;
	}

	@GetMapping("/trace")
	public ResponseEntity<TraceEntrypointResponse> getTraceStep() throws InterruptedException {
		registry.counter("trace_step_total").increment();
		Span span = tracer.currentSpan();
		final String response = traceService.calculateTrace();
		return ResponseEntity.ok(new TraceEntrypointResponse(response));
	}
}
