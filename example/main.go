package main

import (
        "fmt"
        "encoding/json"
        "strings"
        "net/http"

        "github.com/avinoth/go-metacritic/metacritic"
        "github.com/gorilla/mux"
        )

type MovieBasic struct {
  Name, Url, Poster, Certificate, Runtime, ReleaseDate, Genres string
  UserRating Rating
  CriticRating Rating
}

type Movie struct {
  MovieBasic
  CriticReviews []CriticReview
  UserReviews []UserReview
}

type GameBasic struct {
  Name, Url, Summary, ReleaseDate, Certificate, Publisher, Platform string
  CriticRating Rating
}

type Game struct {
  GameBasic
  UserRating Rating
  CriticReviews []CriticReview
  UserReviews []UserReview
}

type AlbumBasic struct {
  Name, Url, Poster, Summary, ReleaseDate, Genres, RecordLabel string
  CriticRating Rating
}

type Album struct {
  AlbumBasic
  UserRating Rating
  CriticReviews []CriticReview
  UserReviews []UserReview
}

type PersonBasic struct {
  Name, Url, DOB, Categories string
  AverageMovieScore Rating
  AverageTVScore Rating
}

type Person struct {
  PersonBasic
  Biography string
  MovieScores Distribution
  TVScores Distribution
  MovieCredits []CreditInfo
  TVCredits []CreditInfo
}

type Distribution struct {
  Positive, Mixed, Negative string
}

type CreditInfo struct {
  CriticRating, Name, Url, Year, Credit, UserRating string
}

type Rating struct {
  Average string
  Count string
}

type CriticReview struct {
  Score, Source, Author, Summary, Url string
}

type UserReview struct {
  Username, ProfileUrl, Score, ReviewDate, Review, Like, Dislike string
}

type TvBasic struct {
  Name, Url, Summary, ReleaseDate, Genres string
  CriticRating Rating
}

func main() {
  r := mux.NewRouter()
  r.HandleFunc("/api/{mode}/{category}/{q}", ServeHandler)
  http.ListenAndServe(":3000", r)
}

func ServeHandler(w http.ResponseWriter, r *http.Request) {
  args := mux.Vars(r)
  category := args["category"]
  query := args["q"]
  mode := args["mode"]
  category = strings.ToLower(category)

  var result string
  var err error

  if (mode == "search") {
    result, err = metacritic.Search(category, query)
  } else if (mode == "find") {
    result, err = metacritic.Find(category, query)
  } else {
    errorHandler(w, r, http.StatusNotFound, "Invalid Mode. Available modes - Search, Find")
  }

  if err != nil {
    fmt.Fprintln(w, err)
  }

  var res []byte
  var movie_results []MovieBasic
  var game_results []GameBasic
  var album_results []AlbumBasic
  var person_results []PersonBasic
  var tv_results []TvBasic
  var movie Movie
  var game Game
  var album Album
  var person Person

  if (mode == "search") {
    switch category {
    case "game":
      err = json.Unmarshal([]byte(result), &game_results)
      if err != nil {
        fmt.Println(err)
        errorHandler(w, r, http.StatusInternalServerError , "")
      }
      res, err = json.Marshal(game_results)
    case "movie":
      err = json.Unmarshal([]byte(result), &movie_results)
      if err != nil {
        fmt.Println(err)
        errorHandler(w, r, http.StatusInternalServerError , "")
      }
      res, err = json.Marshal(movie_results)
    case "album":
      err = json.Unmarshal([]byte(result), &album_results)
      if err != nil {
        fmt.Println(err)
        errorHandler(w, r, http.StatusInternalServerError , "")
      }
      res, err = json.Marshal(album_results)
    case "person":
      err = json.Unmarshal([]byte(result), &person_results)
      if err != nil {
        fmt.Println(err)
        errorHandler(w, r, http.StatusInternalServerError , "")
      }
      res, err = json.Marshal(person_results)
    case "tv":
      err = json.Unmarshal([]byte(result), &tv_results)
      if err != nil {
        fmt.Println(err)
        errorHandler(w, r, http.StatusInternalServerError , "")
      }
      res, err = json.Marshal(tv_results)
    default:
      errorHandler(w, r, http.StatusNotFound, "Invalid Category. Available Game, Movie, Album, Person, Tv")
    }
  } else if (mode == "find") {
    switch category {
      case "game":
        err = json.Unmarshal([]byte(result), &game)
        if err != nil {
          fmt.Println(err)
          errorHandler(w, r, http.StatusInternalServerError , "")
        }
        res, err = json.Marshal(game)
      case "movie":
        err = json.Unmarshal([]byte(result), &movie)
        if err != nil {
          fmt.Println(err)
          errorHandler(w, r, http.StatusInternalServerError , "")
        }
        res, err = json.Marshal(movie)
      case "album":
        err = json.Unmarshal([]byte(result), &album)
        if err != nil {
          fmt.Println(err)
          errorHandler(w, r, http.StatusInternalServerError , "")
        }
        res, err = json.Marshal(album)
      case "person":
        err = json.Unmarshal([]byte(result), &person)
        if err != nil {
          fmt.Println(err)
          errorHandler(w, r, http.StatusInternalServerError , "")
        }
        res, err = json.Marshal(person)
    }
  } else {
    errorHandler(w, r, http.StatusNotFound, "Invalid Mode. Available modes - Search, Find")
  }

  if err != nil {
    fmt.Println(err)
    errorHandler(w, r, http.StatusInternalServerError , "")
  }

  w.Header().Set("Content-Type", "application/json")
  fmt.Fprint(w, string(res))
}

func errorHandler(w http.ResponseWriter, r *http.Request, status int, err_msg string) {
    w.WriteHeader(status)
    fmt.Print(w, err_msg)
}

