<?php
namespace A { class T {} }
namespace B { use A\T as X; echo X::class, "|"; }
namespace C { class X {} echo X::class, "\n"; }
