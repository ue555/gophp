<?php
interface RequiredWork { public function run(); }
echo "before\n";
class IncompleteWork implements RequiredWork {}
echo "after\n";
