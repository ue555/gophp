<?php
echo "before";
class TypedConstantParent { public const int VALUE = 1; }
class TypedConstantInvalidChild extends TypedConstantParent { public const string VALUE = 'bad'; }
echo "after";
