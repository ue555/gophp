<?php
namespace UnionNamespace;
class Base {}
use UnionNamespace\Base as Alias;
echo "before\n";
function wrong(Base|Alias $value) {}
