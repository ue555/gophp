<?php declare(strict_types=1);
interface CrossA {}
interface CrossB {}
function crossTake((CrossA&CrossB)|int $value): (CrossB&CrossA)|int { return $value; }
