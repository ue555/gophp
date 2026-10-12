<?php
// null versus '': == != < <= > >=
echo ((null) == ('')) + 0, ((null) != ('')) + 0, ((null) < ('')) + 0, ((null) <= ('')) + 0, ((null) > ('')) + 0, ((null) >= ('')) + 0, "\n";
// null versus '0': == != < <= > >=
echo ((null) == ('0')) + 0, ((null) != ('0')) + 0, ((null) < ('0')) + 0, ((null) <= ('0')) + 0, ((null) > ('0')) + 0, ((null) >= ('0')) + 0, "\n";
// null versus ' ': == != < <= > >=
echo ((null) == (' ')) + 0, ((null) != (' ')) + 0, ((null) < (' ')) + 0, ((null) <= (' ')) + 0, ((null) > (' ')) + 0, ((null) >= (' ')) + 0, "\n";
// null versus -1: == != < <= > >=
echo ((null) == (-1)) + 0, ((null) != (-1)) + 0, ((null) < (-1)) + 0, ((null) <= (-1)) + 0, ((null) > (-1)) + 0, ((null) >= (-1)) + 0, "\n";
// -1 versus null: == != < <= > >=
echo ((-1) == (null)) + 0, ((-1) != (null)) + 0, ((-1) < (null)) + 0, ((-1) <= (null)) + 0, ((-1) > (null)) + 0, ((-1) >= (null)) + 0, "\n";
// false versus '0': == != < <= > >=
echo ((false) == ('0')) + 0, ((false) != ('0')) + 0, ((false) < ('0')) + 0, ((false) <= ('0')) + 0, ((false) > ('0')) + 0, ((false) >= ('0')) + 0, "\n";
// false versus '': == != < <= > >=
echo ((false) == ('')) + 0, ((false) != ('')) + 0, ((false) < ('')) + 0, ((false) <= ('')) + 0, ((false) > ('')) + 0, ((false) >= ('')) + 0, "\n";
// false versus ' ': == != < <= > >=
echo ((false) == (' ')) + 0, ((false) != (' ')) + 0, ((false) < (' ')) + 0, ((false) <= (' ')) + 0, ((false) > (' ')) + 0, ((false) >= (' ')) + 0, "\n";
// true versus 'false': == != < <= > >=
echo ((true) == ('false')) + 0, ((true) != ('false')) + 0, ((true) < ('false')) + 0, ((true) <= ('false')) + 0, ((true) > ('false')) + 0, ((true) >= ('false')) + 0, "\n";
// true versus -1: == != < <= > >=
echo ((true) == (-1)) + 0, ((true) != (-1)) + 0, ((true) < (-1)) + 0, ((true) <= (-1)) + 0, ((true) > (-1)) + 0, ((true) >= (-1)) + 0, "\n";
// false versus -0.0: == != < <= > >=
echo ((false) == (-0.0)) + 0, ((false) != (-0.0)) + 0, ((false) < (-0.0)) + 0, ((false) <= (-0.0)) + 0, ((false) > (-0.0)) + 0, ((false) >= (-0.0)) + 0, "\n";
// null versus 0.0: == != < <= > >=
echo ((null) == (0.0)) + 0, ((null) != (0.0)) + 0, ((null) < (0.0)) + 0, ((null) <= (0.0)) + 0, ((null) > (0.0)) + 0, ((null) >= (0.0)) + 0, "\n";
// true versus '0': == != < <= > >=
echo ((true) == ('0')) + 0, ((true) != ('0')) + 0, ((true) < ('0')) + 0, ((true) <= ('0')) + 0, ((true) > ('0')) + 0, ((true) >= ('0')) + 0, "\n";
