package org.acesradio.controller;

import io.micronaut.http.HttpHeaders;
import io.micronaut.http.HttpRequest;
import io.micronaut.http.HttpResponse;
import io.micronaut.http.HttpStatus;
import io.micronaut.http.client.HttpClient;
import io.micronaut.http.client.annotation.Client;
import io.micronaut.http.client.exceptions.HttpClientResponseException;
import io.micronaut.runtime.EmbeddedApplication;
import io.micronaut.test.extensions.junit5.annotation.MicronautTest;
import jakarta.inject.Inject;
import org.acesradio.domain.Handle;
import org.acesradio.repository.HandleRepository;
import org.junit.jupiter.api.Test;

import java.util.ArrayList;
import java.util.Collections;
import java.util.List;

import static org.junit.jupiter.api.Assertions.*;


@MicronautTest
class HandleControllerTest {

    @Inject
    @Client("/")
    HttpClient client;

    @Inject
    HandleRepository handleRepository;
    @Inject
    EmbeddedApplication<?> application;

    @Test
    public void testFindNonExistingHandleReturns404() {
        HttpClientResponseException thrown = assertThrows(HttpClientResponseException.class, () -> {
            client.toBlocking().exchange(HttpRequest.GET("/handles/99"));
        });

        assertNotNull(thrown.getResponse());
        assertEquals(HttpStatus.NOT_FOUND, thrown.getStatus());
    }
    @Test
    public void testCreateUser() {

        List<Long> handleIds = new ArrayList<>();

        HttpRequest<?> request = HttpRequest.POST("/handles", Collections.singletonMap("handle", "Drollby Screenface"));
        HttpResponse<?> response = client.toBlocking().exchange(request);
        handleIds.add(entityId(response));

        assertEquals(HttpStatus.CREATED, response.getStatus());

        request = HttpRequest.POST("/handles", Collections.singletonMap("handle", "Crew Eagle"));
        response = client.toBlocking().exchange(request);

        assertEquals(HttpStatus.CREATED, response.getStatus());

    }

    protected Long entityId(HttpResponse<?> response) {
        String path = "/genres/";
        String value = response.header(HttpHeaders.LOCATION);
        if (value == null) {
            return null;
        }
        int index = value.indexOf(path);
        if (index != -1) {
            return Long.valueOf(value.substring(index + path.length()));
        }
        return null;
    }

}



//
//    @Test
//    public void testGenreCrudOperations() {
//
//        List<Long> genreIds = new ArrayList<>();
//
//        HttpRequest<?> request = HttpRequest.POST("/genres", Collections.singletonMap("name", "DevOps"));
//        HttpResponse<?> response = client.toBlocking().exchange(request);
//        genreIds.add(entityId(response));
//
//        assertEquals(HttpStatus.CREATED, response.getStatus());
//
//        request = HttpRequest.POST("/genres", Collections.singletonMap("name", "Microservices"));
//        response = client.toBlocking().exchange(request);
//
//        assertEquals(HttpStatus.CREATED, response.getStatus());
//
//        Long id = entityId(response);
//        genreIds.add(id);
//        request = HttpRequest.GET("/genres/" + id);
//
//        Genre genre = client.toBlocking().retrieve(request, Genre.class);
//
//        assertEquals("Microservices", genre.getName());
//
//        request = HttpRequest.PUT("/genres", new GenreUpdateCommand(id, "Micro-services"));
//        response = client.toBlocking().exchange(request);
//
//        assertEquals(HttpStatus.NO_CONTENT, response.getStatus());
//
//        request = HttpRequest.GET("/genres/" + id);
//        genre = client.toBlocking().retrieve(request, Genre.class);
//        assertEquals("Micro-services", genre.getName());
//
//        request = HttpRequest.GET("/genres/list");
//        List<Genre> genres = client.toBlocking().retrieve(request, Argument.of(List.class, Genre.class));
//
//        assertEquals(2, genres.size());
//
//        request = HttpRequest.POST("/genres/ex", Collections.singletonMap("name", "Microservices"));
//        response = client.toBlocking().exchange(request);
//
//        assertEquals(HttpStatus.NO_CONTENT, response.getStatus());
//
//        request = HttpRequest.GET("/genres/list");
//        genres = client.toBlocking().retrieve(request, Argument.of(List.class, Genre.class));
//
//        assertEquals(2, genres.size());
//
//        request = HttpRequest.GET("/genres/list?size=1");
//        genres = client.toBlocking().retrieve(request, Argument.of(List.class, Genre.class));
//
//        assertEquals(1, genres.size());
//        assertEquals("DevOps", genres.get(0).getName());
//
//        request = HttpRequest.GET("/genres/list?size=1&sort=name,desc");
//        genres = client.toBlocking().retrieve(request, Argument.of(List.class, Genre.class));
//
//        assertEquals(1, genres.size());
//        assertEquals("Micro-services", genres.get(0).getName());
//
//        request = HttpRequest.GET("/genres/list?size=1&page=2");
//        genres = client.toBlocking().retrieve(request, Argument.of(List.class, Genre.class));
//
//        assertEquals(0, genres.size());
//
//        // cleanup:
//        for (Long genreId : genreIds) {
//            request = HttpRequest.DELETE("/genres/" + genreId);
//            response = client.toBlocking().exchange(request);
//            assertEquals(HttpStatus.NO_CONTENT, response.getStatus());
//        }
//    }
//
//    protected Long entityId(HttpResponse<?> response) {
//        String path = "/genres/";
//        String value = response.header(HttpHeaders.LOCATION);
//        if (value == null) {
//            return null;
//        }
//        int index = value.indexOf(path);
//        if (index != -1) {
//            return Long.valueOf(value.substring(index + path.length()));
//        }
//        return null;
//    }