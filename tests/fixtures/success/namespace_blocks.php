<?php
namespace A { const X = 1; function f() { return X; } }
namespace B { const X = 2; function f() { return X; } }
namespace { echo \A\f(), "|", \B\f(), "\n"; }
