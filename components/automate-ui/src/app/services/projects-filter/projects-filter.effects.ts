import { Injectable } from '@angular/core';
import { Router } from '@angular/router';
import { HttpErrorResponse } from '@angular/common/http';
import { of as observableOf } from 'rxjs';
import { switchMap, catchError, map, tap } from 'rxjs/operators';
import { Actions, Effect, ofType } from '@ngrx/effects';
import { ProjectsFilterOption, ProjectsFilterOptionTuple } from './projects-filter.reducer';
import { ProjectsFilterService } from './projects-filter.service';
import {
  ProjectsFilterActionTypes,
  LoadOptions,
  LoadOptionsSuccess,
  LoadOptionsFailure,
  SaveOptions,
  SaveOptionsSuccess,
  SaveOptionsFailure
} from './projects-filter.actions';

@Injectable()
export class ProjectsFilterEffects {
  constructor(
    private actions$: Actions,
    private projectsFilter: ProjectsFilterService,
    private router: Router
  ) { }

  @Effect()
  loadOptions$ = this.actions$.pipe(
    ofType<LoadOptions>(ProjectsFilterActionTypes.LOAD_OPTIONS),
    switchMap(() => {
      return this.projectsFilter.fetchOptions().pipe(
        map((fetched: ProjectsFilterOption[]) => {
          const restored = this.projectsFilter.restoreOptions() || [];
          return new LoadOptionsSuccess(<ProjectsFilterOptionTuple>{
            fetched: fetched,
            restored: restored
          });
        }),
        catchError((error: HttpErrorResponse) => observableOf(new LoadOptionsFailure(error))));
    }));

  @Effect()
  saveOptions$ = this.actions$.pipe(
    ofType<SaveOptions>(ProjectsFilterActionTypes.SAVE_OPTIONS),
    map(({ payload }) => {
      try {
        this.projectsFilter.storeOptions(payload);
        return new SaveOptionsSuccess(payload);
      } catch (error) {
        return new SaveOptionsFailure(error);
      }
    }));

  @Effect({ dispatch: false })
  saveOptionsSuccess$ = this.actions$.pipe(
    ofType<SaveOptionsSuccess>(ProjectsFilterActionTypes.SAVE_OPTIONS_SUCCESS),
    // tap(() => location.reload())
    tap(() => {
      // super hack for avoiding a full location.reload()
      const reloadUrl = '/reload';
      const currentUrl = location.pathname + location.search;
      this.router.navigateByUrl(reloadUrl, { skipLocationChange: true })
        .then(() => this.router.navigateByUrl(currentUrl, { skipLocationChange: true }));
    })
  );

  @Effect({ dispatch: false })
  saveOptionsFailure$ = this.actions$.pipe(
    ofType<SaveOptionsFailure>(ProjectsFilterActionTypes.SAVE_OPTIONS_FAILURE),
    tap(({ payload }) => console.error('Something bad happened: ', payload))
  );
}
