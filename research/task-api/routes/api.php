<?php

use App\Models\Task;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Route;

Route::get('/tasks', fn () => response()->json(Task::orderBy('id')->get()));

Route::post('/tasks', function (Request $request) {
    $data = $request->validate([
        'title' => ['required', 'string', 'max:120'],
        'completed' => ['sometimes', 'required', 'boolean'],
    ]);

    return response()->json(Task::create($data), 201);
});

Route::get('/tasks/{id}', function (int $id) {
    return response()->json(Task::findOrFail($id));
})->whereNumber('id');

Route::patch('/tasks/{id}', function (Request $request, int $id) {
    $task = Task::findOrFail($id);
    $data = $request->validate([
        'title' => ['sometimes', 'required', 'string', 'max:120'],
        'completed' => ['sometimes', 'required', 'boolean'],
    ]);
    $task->fill($data);
    // Empty or unchanged PATCH is a no-op, including updated_at.
    if ($task->isDirty()) {
        $task->save();
    }

    return response()->json($task);
})->whereNumber('id');

Route::delete('/tasks/{id}', function (int $id) {
    $task = Task::findOrFail($id);
    $task->delete();

    return response()->noContent();
})->whereNumber('id');
