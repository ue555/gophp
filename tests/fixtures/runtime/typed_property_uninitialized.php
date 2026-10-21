<?php
class TypedUninitialized { public ?int $value; }
echo "before|";
echo (new TypedUninitialized)->value;
echo "after";
