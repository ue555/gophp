<?php

use Illuminate\Support\Facades\Artisan;

Artisan::command('gophp:probe', function () {
    $this->line('probe-ok');
})->purpose('Check that the fixed application can run a custom command');
