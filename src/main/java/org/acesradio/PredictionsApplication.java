package org.acesradio;

import io.micronaut.runtime.Micronaut;
import io.swagger.v3.oas.annotations.*;
import io.swagger.v3.oas.annotations.info.*;
import org.slf4j.bridge.SLF4JBridgeHandler;

@OpenAPIDefinition(
    info = @Info(
            title = "predictions",
            version = "0.0"
    )
)
public class PredictionsApplication {

    public static void main(String[] args) {
        SLF4JBridgeHandler.removeHandlersForRootLogger();
        SLF4JBridgeHandler.install();

        Micronaut.run(PredictionsApplication.class, args);
    }
}