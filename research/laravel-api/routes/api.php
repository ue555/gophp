<?php

use Illuminate\Http\Request;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Route;

Route::get('/probe', fn () => response()->json(['ok' => true]));

Route::post('/books-manual', function (Request $request) {
    $data = $request->validate(['title' => ['required', 'string', 'max:80']]);
    $id = DB::table('books')->insertGetId([
        'title' => $data['title'],
        'created_at' => now(),
        'updated_at' => now(),
    ]);
    return response()->json(['id' => $id, 'title' => $data['title']], 201);
});

Route::get('/books-manual/{id}', function (int $id) {
    $book = DB::table('books')->where('id', $id)->first();
    return $book === null
        ? response()->json(['error' => 'not_found'], 404)
        : response()->json(['id' => $book->id, 'title' => $book->title]);
});
