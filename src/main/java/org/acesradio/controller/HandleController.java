package org.acesradio.controller;


import io.micronaut.data.model.Pageable;
import io.micronaut.http.HttpResponse;
import io.micronaut.http.HttpStatus;
import io.micronaut.http.annotation.Controller;
import io.micronaut.http.annotation.*;
import io.micronaut.scheduling.TaskExecutors;
import io.micronaut.scheduling.annotation.ExecuteOn;
import org.acesradio.domain.Handle;
import org.acesradio.repository.HandleRepository;

import javax.validation.Valid;
import javax.validation.constraints.NotBlank;
import java.net.URI;
import java.util.List;
import java.util.Optional;

@ExecuteOn(TaskExecutors.IO)
@Controller("/handles")
public class HandleController {

    protected final HandleRepository handleRepository;

    public HandleController(HandleRepository handleRepository) {
        this.handleRepository = handleRepository;
    }

    @Get("/{id}")
    public Optional<Handle> show(Long id) {
        return handleRepository.findById(id);
    }

//    @Put
//    public HttpResponse update(@Body @Valid GenreUpdateCommand command) {
//        genreRepository.update(command.getId(), command.getName());
//        return HttpResponse
//                .noContent()
//                .header(HttpHeaders.LOCATION, location(command.getId()).getPath());
//    }

    @Get("/list")
    public List<Handle> list(@Valid Pageable pageable) {
        return handleRepository.findAll(pageable).getContent();
    }

    @Post
    public HttpResponse<Handle> save(@Body("handle") @NotBlank String handle) {
        Handle theHandle = handleRepository.save(handle);
        return HttpResponse
                .created(theHandle)
                .headers(headers -> headers.location(location(theHandle.id())));
    }


    @Delete("/{id}")
    @Status(HttpStatus.NO_CONTENT)
    public void delete(Long id) {
        handleRepository.deleteById(id);
    }

    protected URI location(Long id) {
        return URI.create("/handles/" + id);
    }

    protected URI location(Handle handle) {
        return location(handle.id());
    }
}