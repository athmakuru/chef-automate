import { Observable } from 'rxjs';
import { first, flatMap } from 'rxjs/operators';
import { Injectable } from '@angular/core';
import { HttpEvent, HttpHandler, HttpInterceptor, HttpRequest } from '@angular/common/http';
import { ProjectsFilterService } from 'app/services/projects-filter/projects-filter.service';
import { ProjectConstants } from 'app/entities/projects/project.model';

@Injectable()
export class ProjectsFilterInterceptor implements HttpInterceptor {
  constructor(private projectsFilter: ProjectsFilterService) {}

  intercept(request: HttpRequest<any>, next: HttpHandler): Observable<HttpEvent<any>> {

    // if (thisEndpointDoesntNeedOrSupportFilteringByProjects) {
    //   // pass the request through
    //   return next.handle(req.clone());
    // }

    return this.projectsFilter.selectionValue$.pipe(
      first(),
      flatMap(selectionValue => {
        let headers = request.headers.set('projects', selectionValue);

        if (selectionValue === ProjectConstants.ALL_RESOURCES_VALUE) {
          headers = headers.delete('projects');
        }

        return next.handle(request.clone({ headers }));
      })
    );
  }
}
