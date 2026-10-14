<?php
namespace Lib { class Item {} function answer() { return 7; } const K = 8; }
namespace App { use Lib\{Item as Box, function answer as value, const K as NUMBER}; echo Box::class, "|", value(), "|", NUMBER, "\n"; }
