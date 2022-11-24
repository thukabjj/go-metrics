package com.example.servicec.service;

import org.springframework.stereotype.Service;

import com.example.servicec.repository.TraceRepository;

@Service
public class TraceService {

	private final TraceRepository repository;

	public TraceService(final TraceRepository repository) {
		this.repository = repository;
	}

	public String calculateTrace() throws InterruptedException {
		Thread.sleep(500);

		final String traceInformation = repository.getTraceInformation();

		return traceInformation;
	}
}
