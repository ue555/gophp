<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Model;

class Task extends Model
{
    protected $fillable = ['title', 'completed'];

    protected $attributes = ['completed' => false];

    protected $casts = ['completed' => 'boolean'];
}
