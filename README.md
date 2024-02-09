# FoodHubber

<img src="https://eventi.emergency.it/wp-content/uploads/cropped-favicon.png" alt="Emergency" width="200"/>
<img src="https://media.licdn.com/dms/image/C4E0BAQFhNBC2FoSLqw/company-logo_200_200/0/1643361232408/aton_spa_logo?e=2147483647&v=beta&t=Z64YPuG9Az_o9LnDX68tmqzAJ_KHMREjg04uk7UjpFY" alt="Aton" width="200"/>

This project is (or, will be) a web system to manage a Food Hub that distributes food and other goods _pro bono_. These food hubs are put in place in the italian territory by Emergency ONG Onlus.

This project is sponsored by [Emergency ONG Onlus](https://emergency.it) and [Aton Società Benefit](https://www.aton.com).


This project is released as Free and Open Source Software under the [GPLv3 license](https://www.gnu.org/licenses/quick-guide-gplv3.it.html).

**NOTE**: The development is currently happening in the `develop` branch. The first versions will be italian only.

## Sources

- [Palette](https://kdesign.co/blog/pastel-color-palette-examples/)

# Installation

1. Create the database;
    1. Create an empty Postgresql engine and connect to it;
    1. As `postgres`, execute the statements in `data/creation.sql`;
    1. Connect to the database `foodhubber` as the user `foodhubber` (password is `foodhubber`);
1. Run the application with the proper parameters (supposing the database is on `localhost:5432`):
```
./foodhubber.exe -db "postgresql://foodhubber:foodhubber@localhost:5432/foodhubber"
```

## Debug

- You can use the `-force-week=1/2/3/4` commandline parameter to simulate being in a week of the month different than the current, real one.